# frozen_string_literal: true

Rails.application.routes.draw do
  devise_for :users, controllers: {
    registrations: "users/registrations",
    sessions: "users/sessions",
    confirmations: "users/confirmations"
  }

  if Rails.env.development?
    mount LetterOpenerWeb::Engine, at: "/letter_opener"
  end

  get "sitemap.xml", to: "sitemap#sitemap", defaults: { format: :xml }
  get "llms.txt",    to: "sitemap#llms",    defaults: { format: :text }

  root "home#index"

  get "up" => "rails/health#show", as: :rails_health_check

  namespace :dashboard do
    root "overview#index"

    resources :api_keys do
      member do
        post :regenerate
        delete :revoke
      end
    end

    resource :usage, only: [ :show ], controller: "usage" do
      collection do
        get :by_endpoint
        get :by_date
        get :export
      end
    end

    resource :billing, only: [ :show, :update ], controller: "billing"
    post "billing/checkout", to: "billing#checkout", as: :checkout_billing
    post "billing/portal", to: "billing#portal", as: :portal_billing
    delete "billing/cancel_subscription", to: "billing#cancel_subscription", as: :cancel_subscription_billing

    resources :invoices, only: [ :index, :show ]

    resource :settings, only: [ :show, :update ] do
      member do
        post   :request_deletion
        get    :confirm_deletion
        delete :execute_deletion
      end
    end
  end

  namespace :admin do
    root "dashboard#index"

    authenticate :user, ->(u) { u.admin? } do
      resources :users do
        member do
          post :suspend
          post :unsuspend
          post :ban
          post :make_admin
          post :remove_admin
        end

        resources :credit_adjustments, only: [ :new, :create ]
      end

      resources :api_keys, only: [ :index, :show ] do
        member do
          delete :revoke
        end
      end

      resources :abuse_reports do
        member do
          post :resolve
          post :investigate
        end
      end

      namespace :analytics do
        get :usage
        get :revenue
        get :system_health
      end
    end
  end

  namespace :webhooks do
    post "lemonsqueezy", to: "lemonsqueezy#create"
  end

  get "docs", to: "home#docs"
  get "pricing", to: "home#pricing"
  get "about", to: "home#about"
  get "team", to: "home#team"
  get "privacy", to: "home#privacy"
  get "terms", to: "home#terms"
  get "contact", to: "home#contact"
  get "api_reference", to: "home#api_reference"
  get "changelog", to: "home#changelog"

  get "blog", to: "home#blog"
  get "status", to: "home#status"
  get "glossary", to: "home#glossary"
  get "error_codes", to: "home#error_codes"
  get "faq", to: "home#faq"

  get "suggest-an-api", to: "suggestions#new", as: "suggest_api"
  post "suggest-an-api", to: "suggestions#create"
  get "talk-to-sales", to: "sales_inquiries#new", as: "talk_to_sales"
  post "talk-to-sales", to: "sales_inquiries#create"

  get "examples", to: "examples#index"
  resources :apis, only: [ :index, :show ]
  resources :categories, only: [ :show ]
  resources :examples, only: [ :show ]

  post "api/proxy", to: "api_proxy#create"
end
